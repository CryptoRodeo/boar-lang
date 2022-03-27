#!/usr/bin/env groovy

pipeline {

    agent any

		stages {
        stage('Build') {
            steps {
                echo 'Building...'
								sh 'go build ./'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing...'
								sh 'go test ./parser'
								sh 'go test ./lexer'
								sh 'go test ./ast'
            }
        }
    }
}