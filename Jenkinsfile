@Library('c1ns-lib-pipeline@v2.3') _

pipeline {
    agent {
        docker 'golang:1.15'
    }

    options {
        buildDiscarder logRotator(artifactNumToKeepStr: '3')
        disableConcurrentBuilds()
    }

    environment {
        RELEASE_BRANCH = 'main'
        OWNER_CHANNEL = 'microwave'
        BUILD_VERSION = buildVersion releaseBranch: RELEASE_BRANCH

        GOCACHE = "${JENKINS_HOME}"
    }

    stages {
        stage('Setup') {
            steps {
                buildName BUILD_VERSION
            }
        }

        stage('Test') {
            steps {
                sh 'make check'
            }
        }

        stage('Build') {
            steps {
                sh 'make build'
            }
        }

        stage('Release') {
            when {
                branch RELEASE_BRANCH
            }
            agent {
                label 'c1ns-build'
            }

            steps {
                gitRelease version: BUILD_VERSION, updateBranch: true
            }
        }
    }

    post {
        unsuccessful {
            script {
                if (env.BRANCH_NAME == RELEASE_BRANCH) {
                    c1netPostBuildFailure ownerChannel: OWNER_CHANNEL
                }
            }
        }

        always {
            script {
                if (env.BRANCH_NAME == RELEASE_BRANCH) {
                    c1netPostBuild()
                } else {
                    c1netPostPRBuild()
                }
            }
        }

        cleanup {
            cleanWs()
        }
    }
}
