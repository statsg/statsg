pipeline {
  agent any
  stages {
    stage('test') {
      steps {
        parallel(
          "test": {
            sh 'make test'
            
          },
          "build": {
            sh 'make build'
            
          }
        )
      }
    }
    stage('deploy') {
      steps {
        echo 'test deploy'
      }
    }
  }
}