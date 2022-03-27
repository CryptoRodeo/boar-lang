#!/usr/bin/env groovy

pipeline {

    agent any

		stages {
        stage('Build') {
            steps {
                echo 'Building...'
								go build .
            }
        }
        stage('Test') {
            steps {
                echo 'Testing...'
								go test ./parser
								go test ./lexer
								go test ./ast
            }
        }
    }
}