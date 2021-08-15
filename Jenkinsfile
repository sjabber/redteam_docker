pipeline {
    agent any

    environment {
      DOCKER_USER_ID = 'sjabber'
      DOCKER_USER_PASSWORD = credentials('dockerPasswd')
    }

    stages {
      stage('Pull') {
        steps {
          echo 'git clone'
          git credentialsId: 'sjabber', url: 'https://github.com/sjabber/redteam_server'
        }
      }
    
      stage('Docker Build') {
        agent any
        steps {
          echo 'build front'
          dir('./redteam') {
            sh(script: '''
              pwd
              docker build -f redteam.Dockerfile -t sjabber/redteam .
              docker build -f redteam2.Dockerfile -t sjabber/redteam_java .
              docker build -f redteam_front.Dockerfile -t sjabber/redteam_front .
              '''
            )
          }
        }
      }

      stage('Tag') {
        agent any

        steps {
          sh(script: '''
          docker tag sjabber/redteam_front \
          sjabber/redteam_front:${BUILD_NUMBER}

          docker tag sjabber/redteam \
          sjabber/redteam:${BUILD_NUMBER}

          docker tag sjabber/redteam_java \
          sjabber/redteam_java:${BUILD_NUMBER}
          ''')
        }
      }

      stage('Docker Push') {
        agent any

        steps {
          sh(script: 'docker login -u ${DOCKER_USER_ID} -p ${DOCKER_USER_PASSWORD}')
          sh(script: 'docker push ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}')
          sh(script: 'docker push ${DOCKER_USER_ID}/redteam_java:${BUILD_NUMBER}')
          sh(script: 'docker push ${DOCKER_USER_ID}/redteam:${BUILD_NUMBER}')
        }
      }

      stage('AWS Deploy') {
        agent any

        // 기존에 작동중이던 도커 컨테이너 중지, 삭제
        steps {
          echo "redteam deploy start"
          sshagent(['ec2-server']) {
            sh "ssh -o StrictHostKeyChecking=no ubuntu@15.165.17.133 sudo docker rm -f redteam_front"
            sh "ssh -o StrictHostKeyChecking=no ubuntu@15.165.17.133 sudo docker rm -f redteam_java"
            sh "ssh -o StrictHostKeyChecking=no ubuntu@15.165.17.133 sudo docker rm -f redteam"
          }

        }

        post {
          // 기존에 작동중이던 컨테이너가 있는 경우
          always {
            sshagent(['ec2-server']) {
              sh "ssh -o StrictHostKeyChecking=no ubuntu@15.165.17.133 sudo docker login -u ${DOCKER_USER_ID} -p ${DOCKER_USER_PASSWORD}"
              sh "ssh -o StrictHostKeyChecking=no ubuntu@15.165.17.133 sudo docker run -d -p 80:80 --name redteam_front \
            ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}"
              sh "ssh -o StrictHostKeyChecking=no ubuntu@15.165.17.133 sudo docker run -d -p 5001:5001 --name redteam_java \
            ${DOCKER_USER_ID}/redteam_java:${BUILD_NUMBER}"
              sh "ssh -o StrictHostKeyChecking=no ubuntu@15.165.17.133 sudo docker run -d -p 5000:5000 --name redteam \
            ${DOCKER_USER_ID}/redteam:${BUILD_NUMBER}"
            }

          }
        }
      }
    }
}