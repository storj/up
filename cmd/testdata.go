// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
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
	"storj.io/storj/satellite"
	"storj.io/storj/satellite/accounting"
	"storj.io/storj/satellite/buckets"
	"storj.io/storj/satellite/compensation"
	"storj.io/storj/satellite/nodeselection"
	"storj.io/storj/satellite/satellitedb"
)

var (
	database, email, bucket, useragent, period string
	gb                                         = decimal.NewFromInt(1e9)
	tb                                         = decimal.NewFromInt(1e12)
	getRate                                    = int64(20)
	auditRate                                  = int64(10)
	storageRate                                = 0.00000205
)

var testdataCmd = &cobra.Command{
	Use:   "testdata",
	Args:  cobra.NoArgs,
	Short: "Generate testdata to the database",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

func paymentCmd() *cobra.Command {
	paymentCmd := &cobra.Command{
		Use:   "payment",
		Short: "Generate payment and paystub entries for each node",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return generatePayments(database)
		},
	}
	paymentCmd.PersistentFlags().StringVarP(&database, "database", "d", "cockroach://root@localhost:26257/master?sslmode=disable", "Database connection string to generate data")
	return paymentCmd
}

func projectUsageCmd() *cobra.Command {
	projectUsageCmd := &cobra.Command{
		Use:   "project-usage",
		Short: "Generated bandwidth rollups for buckets and projects",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if period == "" {
				period = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1).Format("2006-01")
			}
			usagePeriod, err := time.Parse("2006-01", period)
			if err != nil {
				return errs.New("invalid date specified specified. accepted format is yyyy-mm: %v", err)
			}
			return generateProjectUsage(database, email, bucket, useragent, usagePeriod)
		},
	}
	projectUsageCmd.PersistentFlags().StringVarP(&database, "database", "d", "spanner://projects/test-project/instances/test-instance/databases/master", "Database connection string to generate data")
	projectUsageCmd.PersistentFlags().StringVarP(&email, "email", "e", "test@storj.io", "the email address of the user to add data for")
	projectUsageCmd.PersistentFlags().StringVarP(&bucket, "bucket", "b", "storage-bucket", "the bucket to add the usage for")
	projectUsageCmd.PersistentFlags().StringVarP(&useragent, "useragent", "u", "", "useragent for value attribution")
	projectUsageCmd.PersistentFlags().StringVarP(&period, "period", "p", "", "the month to add usage for. defaults to the previous month")

	return projectUsageCmd
}

func init() {
	RootCmd.AddCommand(testdataCmd)
	testdataCmd.AddCommand(paymentCmd())
	testdataCmd.AddCommand(projectUsageCmd())
}

func generateProjectUsage(database, email string, bucketname string, useragent string, period time.Time) error {
	ctx := context.Background()
	db, err := satellitedb.Open(ctx, zap.L().Named("db"), database, satellitedb.Options{ApplicationName: "satellite-compensation"})
	if err != nil {
		return errs.Wrap(err)
	}
	defer func() {
		_ = db.Close()
	}()

	projects, err := db.Console().Projects().GetAll(ctx)
	if err != nil {
		return errs.Wrap(err)
	}

	for _, p := range projects {

		firstDayOfMonth := time.Date(period.Year(), period.Month(), 1, 1, 0, 0, 0, period.Location())
		lastDayOfMonth := time.Date(period.Year(), period.Month(), 1, 0, 0, 0, 0, period.Location()).AddDate(0, 1, -1)

		var bucket buckets.Bucket
		bucket, err = db.Buckets().GetBucket(ctx, []byte(bucketname), p.ID)
		if err != nil {
			if buckets.ErrBucketNotFound.Has(err) {
				var bucketID uuid.UUID
				bucketID, err = uuid.New()
				if err != nil {
					return errs.Wrap(err)
				}
				// try to create it instead
				bucket, err = db.Buckets().CreateBucket(ctx, buckets.Bucket{
					ID:                          bucketID,
					Name:                        bucketname,
					ProjectID:                   p.ID,
					UserAgent:                   []byte(useragent),
					Created:                     firstDayOfMonth,
					PathCipher:                  0,
					DefaultRedundancyScheme:     storj.RedundancyScheme{},
					DefaultEncryptionParameters: storj.EncryptionParameters{},
					Placement:                   0,
				})
			}
			if err != nil {
				// couldn't get nor create bucket
				return errs.Wrap(err)
			}
		}

		StoredData := int64(1583717400000)
		MetadataSize := int64(2)
		Object := int64(1)
		SegmentCount := int64(2)

		err = updateUsage(ctx, db, crateTally(bucket.Name, p.ID, firstDayOfMonth, Object, SegmentCount, StoredData, MetadataSize))
		if err != nil {
			return errs.Wrap(err)
		}
		err = updateUsage(ctx, db, crateTally(bucket.Name, p.ID, firstDayOfMonth.Add(1*time.Minute), Object, SegmentCount, StoredData, MetadataSize))
		if err != nil {
			return errs.Wrap(err)
		}
		err = updateUsage(ctx, db, crateTally(bucket.Name, p.ID, lastDayOfMonth, Object, SegmentCount, StoredData, MetadataSize))
		if err != nil {
			return errs.Wrap(err)
		}
		err = updateUsage(ctx, db, crateTally(bucket.Name, p.ID, lastDayOfMonth.Add(1*time.Minute), Object, SegmentCount, StoredData, MetadataSize))
		if err != nil {
			return errs.Wrap(err)
		}
		intervalStart := firstDayOfMonth
		for i := 0; i < 24; i++ {
			usage := 1024000000000
			err = db.Orders().UpdateBucketBandwidthAllocation(ctx, p.ID, []byte(bucket.Name), pb.PieceAction_GET, int64(usage), intervalStart)
			if err != nil {
				return errs.Wrap(err)
			}
			err = db.Orders().UpdateBucketBandwidthSettle(ctx, p.ID, []byte(bucket.Name), pb.PieceAction_GET, int64(usage), 0, intervalStart)
			if err != nil {
				return errs.Wrap(err)
			}
			intervalStart = intervalStart.Add(1 * time.Hour)
		}

		createdAt := time.Date(period.Year(), period.Month()-1, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		_, err = db.Testing().RawDB().Exec(ctx, "UPDATE stripe_customers SET created_at = $1 WHERE true", createdAt)
		if err != nil {
			return errs.Wrap(err)
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

		err = db.OverlayCache().IterateAllContactedNodes(ctx, func(ctx context.Context, node *nodeselection.SelectedNode) error {
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
			return errs.Wrap(err)
		}
	}
	err = db.Compensation().RecordPeriod(ctx, paystubs, payments)
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func crateTally(bucketName string, projectID uuid.UUID, intervalStart time.Time, objectCount int64, totalSegmentCount int64,
	totalBytes int64, metadataSize int64) accounting.BucketStorageTally {
	return accounting.BucketStorageTally{
		BucketName:        bucketName,
		ProjectID:         projectID,
		IntervalStart:     intervalStart,
		ObjectCount:       objectCount,
		TotalSegmentCount: totalSegmentCount,
		TotalBytes:        totalBytes,
		MetadataSize:      metadataSize,
	}
}

func updateUsage(ctx context.Context, db satellite.DB, tally accounting.BucketStorageTally) error {
	query := `UPDATE bucket_storage_tallies SET
		total_bytes = bucket_storage_tallies.total_bytes + ?,
		inline = bucket_storage_tallies.inline + ?,
		remote = bucket_storage_tallies.remote + ?,
		total_segments_count = bucket_storage_tallies.total_segments_count + ?,
		remote_segments_count = bucket_storage_tallies.remote_segments_count + ?,
		inline_segments_count = bucket_storage_tallies.inline_segments_count + ?,
		object_count = bucket_storage_tallies.object_count + ?,
		metadata_size = bucket_storage_tallies.metadata_size + ?
		WHERE interval_start=?
		AND bucket_name=?
		AND project_id=?;`
	result, err := db.Testing().RawDB().Exec(ctx,
		query,
		tally.TotalBytes, 0, 0,
		tally.TotalSegmentCount, 0, 0,
		tally.ObjectCount, tally.MetadataSize,
		tally.IntervalStart,
		[]byte(tally.BucketName), tally.ProjectID)

	if err != nil {
		return errs.Wrap(err)
	}

	numRows, _ := result.RowsAffected()
	if numRows < 1 {
		query := `INSERT INTO bucket_storage_tallies (
			interval_start,
			bucket_name, project_id,
			total_bytes, inline, remote,
			total_segments_count, remote_segments_count, inline_segments_count,
			object_count, metadata_size)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		_, err := db.Testing().RawDB().Exec(ctx,
			query,
			tally.IntervalStart,
			[]byte(tally.BucketName), tally.ProjectID,
			tally.TotalBytes, 0, 0,
			tally.TotalSegmentCount, 0, 0,
			tally.ObjectCount, tally.MetadataSize)
		return errs.Wrap(err)
	}
	return errs.Wrap(err)
}
