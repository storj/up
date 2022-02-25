pipeline {
    agent {
        docker {
            label 'main'
            image 'storjlabs/ci:latest'
            alwaysPull true
            args '-u root:root --cap-add SYS_PTRACE -v "/tmp/gomod":/go/pkg/mod'
        }
    }

    options {
          timeout(time: 26, unit: 'MINUTES')
    }

    environment {
        GOTRACEBACK = 'all'
    }

    stages {
        stage('Checkout') {
            steps {
               // Delete any content left over from a previous run.
               sh "chmod -R 777 ."

               // Bash requires extglob option to support !(.git) syntax,
               // and we don't want to delete .git to have faster clones.
               sh 'bash -O extglob -c "rm -rf !(.git)"'

               checkout scm
            }
        }
		stage('Build') {
			parallel {
				stage('Lint') {
					steps {
						sh "go install github.com/magefile/mage@v1.11.0"
						sh "mage -v lint"
					}
				}
				stage('Test') {
					steps {
						sh "go install github.com/magefile/mage@v1.11.0"
						sh "mage -v test"
					}
				}
			}
		}
    }
}
