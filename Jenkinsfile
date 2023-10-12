pipeline {
    agent {
       label('nomad || ondemand')
    }

    options {
          timeout(time: 26, unit: 'MINUTES')
    }

    environment {
        GOTRACEBACK = 'all'
        NO_COLOR = '1'
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
            parallel {
                stage('Uplink') {
                    steps {
                        sh "earthly -P ./test/uplink+integration-uplink"
                    }
                }
                stage('Edge') {
                    steps {
                        sh "earthly -P ./test/edge+integration-edge"
                    }
                }
                stage('Storjscan') {
                    steps {
                        sh "earthly -P ./test/storjscan+integration-storjscan"
                    }
                }
            }
        }
    }
}
