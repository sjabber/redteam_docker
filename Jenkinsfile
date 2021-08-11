pipeline {
    agent any //아무 노예나 써라.

    triggers {
        pollSCM('*/3 * * * *') // cron syntax, 3분 주기로 파이프라인 구동하는 트리거
    }

    stages {
      stage('Pull') {
        steps {
          echo 'git clone'
          git credentialsId: 'sjabber', url: 'https://github.com/sjabber/redteam_server'
        }
    }
    
    stage('Build') {
      agent any
      steps {
        echo 'build front'
        dir('./redteam') {
          sh(script: '''
            pwd
            docker build -f redteam_front.Dockerfile -t sjabber/redteam_front .
            '''
          )
        }
      }
    }
    stage('Tag') {
      agent any

      steps {
        sh(script: '''docker tag sjabber/redteam_front \
        ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}}''')
      }
    }
    stage('Push') {
      agent any

      steps {
        sh(script: 'docker login -u ${DOKER_USER_ID} -p ${DOCKER_USER_PASSWORD}')
        sh(script: 'docker push ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}')
        sh(script: 'docker push ${DOCKER_USER_ID}/redteam_front:latest')
      }
    }
    stage('Deploy') {
      agent any

      steps {
        echo "redteam deploy start"
        sh(script: 'docker stop redteam')
        sh(script: 'docker rm redteam')
        sh(script: '''docker run -d -p 80:80 --name redteam_front \
        ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}''')
      }
    }
  }
}



// node {

//   withCredentials([[$class: 'UsernamePasswordMultiBinding',
//   credentialsId: 'dockerhub',
//   usernameVariable: 'Docker_USER_ID',
//   passwordVariable: 'DOCKER_USER_PASSWORD']]) {
//     stage('Pull') {
//         git credentialsId: 'sjabber', url: 'https://github.com/sjabber/redteam_server'
//     }
//     stage('Build') {
//         sh(script: 'docker build --force-rm=true -f redteam_front.Dockerfile -t ${DOKER_USER_ID}/redteam_front ./redteam')
//     }
//     stage('Tag') {
//         sh(script: '''docker tag ${DOKER_USER_ID}/redteam_front \
//         ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}}''')
//     }
//     stage('Push') {
//       sh(script: 'docker login -u ${DOKER_USER_ID} -p ${DOCKER_USER_PASSWORD}')
//       sh(script: 'docker push ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}')
//       sh(script: 'docker push ${DOCKER_USER_ID}/redteam_front:latest')
//     }
//     stage('Deploy') {
//       try{
//         echo "redteam deploy start"
//         sh(script: 'docker stop redteam')
//         sh(script: 'docker rm redteam')
//       } catch(e) {
//         echo "No redteam container exists"
//       }
//       sh(script: '''docker run -d -p 80:80 --name redteam_front \
//       ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}''')
//     }
//   }
// }