#!/usr/bin/env groovy

pipeline {
    agent any

	stages {
        stage('Build') {
            steps {
                echo 'Building...'
                sh '    ./misc/build.sh'
            }
        }

        stage('Test') {
            steps {
                echo 'Testing...'
                sh '    ./misc/test.sh'
            }
        }
    }
}
