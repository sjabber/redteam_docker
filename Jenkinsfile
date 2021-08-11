node {

  withCredentials([[$class: 'UsernamePasswordMultiBinding',
  credentialsId: 'dockerhub',
  usernameVariable: 'Docker_USER_ID',
  passwordVariable: 'DOCKER_USER_PASSWORD']]) {
    stage('Pull') {
        git credentialsId: 'sjabber', url: 'https://github.com/sjabber/redteam_server'
    }
     stage('Initialize'){
        def dockerHome = tool 'myDocker'
        env.PATH = "${dockerHome}/bin:${env.PATH}"
    }
    stage('Unit Test') {

    }
    stage('Build') {
        sh(script: '''
        #!/bin/bash
        cd ./redteam
        docker build -f redteam_front.Dockerfile -t sjabber/redteam_front .
        ''')
    }
    stage('Tag') {
        sh(script: '''docker tag ${DOKER_USER_ID}/redteam_front \
        ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}}''')
    }
    stage('Push') {
      sh(script: 'docker login -u ${DOKER_USER_ID} -p ${DOCKER_USER_PASSWORD}')
      sh(script: 'docker push ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}')
      sh(script: 'docker push ${DOCKER_USER_ID}/redteam_front:latest')
    }
    stage('Deploy') {
      try{
        echo "redteam deploy start"
        sh(script: 'docker stop redteam')
        sh(script: 'docker rm redteam')
      } catch(e) {
        echo "No redteam container exists"
      }
      sh(script: '''docker run -d -p 80:80 --name redteam_front \
      ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}''')
    }
  }
}