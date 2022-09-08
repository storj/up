pipeline {
    agent {
       label 'node4'
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
               checkout scm
            }
        }
	stage('Lint') {
		steps {
			sh "earthly +lint"
		}
	}
	stage('Test') {
		steps {
			sh "earthly +test"
		}
	}
	stage('Integration') {
		steps {
			sh "earthly -P +integration"
		}
	}
    }
}
