#!/usr/bin/env groovy

pipeline {
    agent any

	stages {
        stage('Build') {
            steps {
                echo 'Building...'
                sh 'bash ./misc/build.sh'
            }
        }

        stage('Test') {
            steps {
                echo 'Testing...'
                sh 'bash ./misc/test.sh'
            }
        }
    }
}
