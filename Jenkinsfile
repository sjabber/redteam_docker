pipeline {
    agent any

    triggers {
        pollSCM('*/1 * * * *') // cron syntax, 1분 주기로 파이프라인 구동하는 트리거
    }

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
    
      stage('Build') {
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

      stage('Push') {
        agent any

        steps {
          sh(script: 'docker login -u ${DOCKER_USER_ID} -p ${DOCKER_USER_PASSWORD}')
          sh(script: 'docker push ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}')
          sh(script: 'docker push ${DOCKER_USER_ID}/redteam_java:${BUILD_NUMBER}')
          sh(script: 'docker push ${DOCKER_USER_ID}/redteam:${BUILD_NUMBER}')
        }
      }

      stage('Deploy') {
        agent any

        // def dockerRun1 = 'docker run -d -p 80:80 --name redteam_front \
        // ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}'

        // def dockerRun2 = 'docker run -d -p 5001:5001 --name redteam_java \
        // ${DOCKER_USER_ID}/redteam_java:${BUILD_NUMBER}'

        // def dockerRun3 = 'docker run -d -p 5000:5000 --name redteam \
        // ${DOCKER_USER_ID}/redteam:${BUILD_NUMBER}'

        // 기존에 작동중이던 도커 컨테이너 중지, 삭제
        steps {
          echo "redteam deploy start"
          sshagent(['ec2-server']) {
            sh "ssh -o StrictHostKeyChecking=no ubuntu@15.165.17.133 sudo docker rm -f `docker ps -q -a`"
          }


          // sh(script: '''
          // docker stop redteam_front
          // docker stop redteam_java
          // docker stop redteam
          // docker rm redteam_front
          // docker rm redteam_java
          // docker rm redteam
          // ''')
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

            // sh(script: '''
            // docker run -d -p 80:80 --name redteam_front \
            // ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}

            // docker run -d -p 5001:5001 --name redteam_java \
            // ${DOCKER_USER_ID}/redteam_java:${BUILD_NUMBER}

            // docker run -d -p 5000:5000 --name redteam \
            // ${DOCKER_USER_ID}/redteam:${BUILD_NUMBER}
            // ''')
          }

          // 기존에 작동중이던 컨테이너가 없는 경우
          // failure {
          //   sshagent(['ec2-server']) {
          //     sh "ssh -o StrictHostKeyChecking=no ubuntu@15.165.17.133 sudo docker run -d -p 80:80 --name redteam_front \
          //   ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}"
          //     sh "ssh -o StrictHostKeyChecking=no ubuntu@15.165.17.133 sudo docker run -d -p 5001:5001 --name redteam_java \
          //   ${DOCKER_USER_ID}/redteam_java:${BUILD_NUMBER}"
          //     sh "ssh -o StrictHostKeyChecking=no ubuntu@15.165.17.133 sudo docker run -d -p 5000:5000 --name redteam \
          //   ${DOCKER_USER_ID}/redteam:${BUILD_NUMBER}"
          //   }

            // sh(script: '''
            // docker run -d -p 80:80 --name redteam_front \
            // ${DOCKER_USER_ID}/redteam_front:${BUILD_NUMBER}

            // docker run -d -p 5001:5001 --name redteam_java \
            // ${DOCKER_USER_ID}/redteam_java:${BUILD_NUMBER}

            // docker run -d -p 5000:5000 --name redteam \
            // ${DOCKER_USER_ID}/redteam:${BUILD_NUMBER}
            // ''')
          //}
        }
      }

      // stage('delete oldversion') {
      //   agent any
      //   steps {
      //     echo "delete old docker images"
      //     sh(script: '''
      //     #! /bin/bash
      //     newversion = ${BUILD_NUMBER}
      //     oldversion = newversion - 1
      //     docker rmi ${DOCKER_USER_ID}/redteam_front:$oldversion
      //     docker rmi ${DOCKER_USER_ID}/redteam_java:$oldversion
      //     docker rmi ${DOCKER_USER_ID}/redteam:$oldversion 
      //     ''')
      //   }
      // }
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