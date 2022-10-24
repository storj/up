// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs"
	"go.uber.org/zap"

	"storj.io/common/pb"
	"storj.io/common/storj"
	"storj.io/common/uuid"
	"storj.io/storj/private/currency"
	"storj.io/storj/satellite/accounting"
	"storj.io/storj/satellite/compensation"
	"storj.io/storj/satellite/overlay"
	"storj.io/storj/satellite/satellitedb"
)

var (
	gb          = decimal.NewFromInt(1e9)
	tb          = decimal.NewFromInt(1e12)
	getRate     = int64(20)
	auditRate   = int64(10)
	storageRate = 0.00000205
)

func testdataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testdata",
		Short: "Generate testdata to the database",
	}

	database := cmd.Flags().String("database", "cockroach://root@localhost:26257/master?sslmode=disable", "Database connection string to generate data")
	var generators []*cobra.Command
	{
		subCmd := &cobra.Command{
			Use:   "payment",
			Short: "Generate payment and paystub entries for each node",
			RunE: func(cmd *cobra.Command, args []string) error {
				return generatePayments(*database)
			},
		}
		generators = append(generators, subCmd)
	}

	{
		subCmd := &cobra.Command{
			Use:   "project-usage",
			Short: "Generated bandwidth rollups for buckets and projects",
			RunE: func(cmd *cobra.Command, args []string) error {
				return generateProjectUsage(*database)
			},
		}
		generators = append(generators, subCmd)
	}

	{
		subCmd := &cobra.Command{
			Use:   "fix-billing",
			Short: "fixes billing, by overriding created_at timestamp",
			RunE: func(cmd *cobra.Command, args []string) error {
				return fixBilling()
			},
		}
		generators = append(generators, subCmd)
	}

	{
		subCmd := &cobra.Command{
			Use:   "all",
			Short: "Execute all the data generators",
			RunE: func(cmd *cobra.Command, args []string) error {
				for _, g := range generators {
					err := g.RunE(cmd, args)
					if err != nil {
						zap.L().Error("Couldn't execute generator", zap.Error(err))
					}
				}
				return nil
			},
		}
		cmd.AddCommand(subCmd)
	}

	cmd.AddCommand(generators...)
	return cmd
}

func fixBilling() error {
	db, err := sql.Open("pgx", "host=localhost port=26257 user=root dbname=master sslmode=disable")
	if err != nil {
		return errs.Wrap(err)
	}
	defer func() {
		_ = db.Close()
	}()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	sqlStatement := `
	UPDATE stripe_customers
	SET created_at ='2020-05-10'
	WHERE true;`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
	zap.L().Error("Couldn't execute generator", zap.Error(err))
	fmt.Print("\nfixed created-at date, invoice command should work now\n")

	return nil
}

func generateProjectUsage(database string) error {
	ctx := context.Background()
	db, err := satellitedb.Open(ctx, zap.L().Named("db"), database, satellitedb.Options{ApplicationName: "satellite-compensation"})
	if err != nil {
		return errs.Wrap(err)
	}
	defer func() {
		_ = db.Close()
	}()

	byEmail, err := db.Console().Users().GetByEmail(ctx, "test@storj.io")
	if err != nil {
		return err
	}

	projects, err := db.Console().Projects().GetAll(ctx)
	if err != nil {
		return err
	}

	for _, p := range projects {

		now := time.Now()
		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()
		lastOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation).AddDate(0, 0, -1)
		dayTen := time.Date(currentYear, currentMonth, 10, 1, 0, 0, 0, currentLocation).AddDate(0, -1, 0)
		intervalStart := dayTen

		bucket, err := db.Buckets().CreateBucket(ctx, storj.Bucket{
			ID:                          byEmail.ID,
			Name:                        "storage-bucket",
			ProjectID:                   p.ID,
			PartnerID:                   uuid.UUID{},
			UserAgent:                   nil,
			Created:                     dayTen,
			PathCipher:                  0,
			DefaultRedundancyScheme:     storj.RedundancyScheme{},
			DefaultEncryptionParameters: storj.EncryptionParameters{},
			Placement:                   0,
		})
		if err != nil {
			fmt.Printf("Unable to create bucket: %s\n", err.Error())
		}

		StoredData := int64(1583717400000)
		MetadataSize := int64(2)
		Object := int64(1)
		SegmentCount := int64(2)

		tally := accounting.BucketStorageTally{
			BucketName:        bucket.Name,
			ProjectID:         p.ID,
			IntervalStart:     dayTen,
			ObjectCount:       Object,
			TotalSegmentCount: SegmentCount,
			TotalBytes:        StoredData,
			MetadataSize:      MetadataSize,
		}

		err = db.ProjectAccounting().CreateStorageTally(ctx, tally)
		if err != nil {
			return err
		}
		tally = accounting.BucketStorageTally{
			BucketName:        bucket.Name,
			ProjectID:         p.ID,
			IntervalStart:     dayTen.Add(1 * time.Minute),
			ObjectCount:       Object,
			TotalSegmentCount: SegmentCount,
			TotalBytes:        StoredData,
			MetadataSize:      MetadataSize,
		}
		err = db.ProjectAccounting().CreateStorageTally(ctx, tally)
		if err != nil {
			return err
		}

		tally = accounting.BucketStorageTally{
			BucketName:        bucket.Name,
			ProjectID:         p.ID,
			IntervalStart:     lastOfMonth,
			ObjectCount:       Object,
			TotalSegmentCount: SegmentCount,
			TotalBytes:        StoredData,
			MetadataSize:      MetadataSize,
		}

		err = db.ProjectAccounting().CreateStorageTally(ctx, tally)
		if err != nil {
			return err
		}
		tally = accounting.BucketStorageTally{
			BucketName:        bucket.Name,
			ProjectID:         p.ID,
			IntervalStart:     lastOfMonth.Add(1 * time.Minute),
			ObjectCount:       Object,
			TotalSegmentCount: SegmentCount,
			TotalBytes:        StoredData,
			MetadataSize:      MetadataSize,
		}
		err = db.ProjectAccounting().CreateStorageTally(ctx, tally)
		if err != nil {
			return err
		}
		for i := 0; i < 24; i++ {
			usage := 1024000000000
			err = db.Orders().UpdateBucketBandwidthAllocation(ctx, p.ID, []byte(bucket.Name), pb.PieceAction_GET, int64(usage), intervalStart)
			if err != nil {
				return err
			}
			err = db.Orders().UpdateBucketBandwidthSettle(ctx, p.ID, []byte(bucket.Name), pb.PieceAction_GET, int64(usage), 0, intervalStart)
			if err != nil {
				return err
			}
			intervalStart = intervalStart.Add(-1 * time.Hour)
		}

	}
	return nil
}

func generatePayments(database string) error {
	ctx := context.Background()
	db, err := satellitedb.Open(ctx, zap.L().Named("db"), database, satellitedb.Options{ApplicationName: "satellite-compensation"})
	if err != nil {
		return errs.Wrap(err)
	}
	defer func() {
		_ = db.Close()
	}()

	db.StoragenodeAccounting()
	var paystubs []compensation.Paystub
	var payments []compensation.Payment
	now := time.Now()
	paymentTypes := []string{"eth", "zksync", "polygon"}
	for i := 0; i < 10; i++ {
		oneMonthBefore := now.AddDate(0, -i, 0)
		period := compensation.Period{
			Year:  oneMonthBefore.Year(),
			Month: oneMonthBefore.Month(),
		}

		err = db.OverlayCache().IterateAllContactedNodes(ctx, func(ctx context.Context, node *overlay.SelectedNode) error {
			storedDataGB := rand.Intn(1000) + 400
			getUsage := int64(storedDataGB * 10 / 7)
			paystub := compensation.Paystub{
				Period:         period,
				NodeID:         compensation.NodeID(node.ID),
				UsageAtRest:    float64(storedDataGB * 24 * 30),
				UsageGet:       getUsage,
				UsagePut:       getUsage * 11 / 10,
				UsageGetAudit:  getUsage / 800000,
				UsageGetRepair: getUsage / 2500,
				UsagePutRepair: getUsage / 30,
			}

			paystub.CompAtRest, err = currency.MicroUnitFromDecimal(
				decimal.NewFromFloat(paystub.UsageAtRest).
					Mul(decimal.NewFromFloat(storageRate)).
					Div(gb))
			if err != nil {
				return errs.Wrap(err)
			}

			paystub.CompGet, err = currency.MicroUnitFromDecimal(
				decimal.NewFromInt(paystub.UsageGet).
					Mul(decimal.NewFromInt(getRate)).
					Div(tb))
			if err != nil {
				return errs.Wrap(err)
			}

			paystub.CompGetAudit, err = currency.MicroUnitFromDecimal(
				decimal.NewFromInt(paystub.UsageGetAudit).
					Mul(decimal.NewFromInt(auditRate)).
					Div(tb))
			if err != nil {
				return errs.Wrap(err)
			}

			paystub.CompGetRepair, err = currency.MicroUnitFromDecimal(
				decimal.NewFromInt(paystub.UsagePutRepair).
					Mul(decimal.NewFromInt(auditRate)).
					Div(tb))
			if err != nil {
				return errs.Wrap(err)
			}

			paystub.CompPutRepair, err = currency.MicroUnitFromDecimal(
				decimal.NewFromInt(paystub.UsageGetRepair).
					Mul(decimal.NewFromInt(auditRate)).
					Div(tb))
			if err != nil {
				return errs.Wrap(err)
			}

			paystub.Paid, err = currency.MicroUnitFromDecimal(
				paystub.CompAtRest.Decimal().Add(
					paystub.CompGet.Decimal()).Add(
					paystub.CompGetAudit.Decimal()).Add(
					paystub.CompGetRepair.Decimal()).Add(
					paystub.CompPutRepair.Decimal()))
			if err != nil {
				return errs.Wrap(err)
			}

			paystubs = append(paystubs, paystub)
			receipt := paymentTypes[i%3] + ":0xc6d9062f010b8c1efd37e65851cc55d4c258af7df2425f766ca9aab4b2b26360"
			payments = append(payments, compensation.Payment{
				Period:  period,
				NodeID:  compensation.NodeID(node.ID),
				Amount:  currency.NewMicroUnit(rand.Int63n(10000) + 10000),
				Receipt: &receipt,
			})
			return nil
		})
		if err != nil {
			return err
		}
	}
	err = db.Compensation().RecordPeriod(ctx, paystubs, payments)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RootCmd.AddCommand(testdataCmd())
}
