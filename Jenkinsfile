#!/usr/bin/env groovy

pipeline {
    agent any

		stages {
        stage('Build') {
            steps {
                echo 'Building...'
								sh '/usr/local/go/bin/go build ./'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing...'
								sh '/usr/local/go/bin/go test ./parser'
								sh '/usr/local/go/bin/go test ./lexer'
								sh '/usr/local/go/bin/go test ./ast'
            }
        }
    }
}
