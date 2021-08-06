pipeline {
    // 스테이지 별로 다른 거
    agent any //agent 무관

    triggers {
        pollSCM('*/3 * * * *') // cron syntax, 3분 주기로 작동하는 파이프라인 트리거
    }

    environment { // 파이프라인 안에서 쓸 환경변수를 입혀준다.
      AWS_ACCESS_KEY_ID = credentials('awsAccessKey') //aws cli 엑세스 키
      AWS_SECRET_ACCESS_KEY = credentials('awsSecretAccessKey')
      AWS_DEFAULT_REGION = 'ap-northeast-2' // 서울을 기준으로 삼아준다.
      HOME = '.' // Avoid npm root owned
    }

    stages {
        // 레포지토리를 다운로드 받음
        stage('Pull') {
            agent any
            
            steps {
                echo 'Clonning Repository'

                git 'https://github.com/sjabber/redteam_server.git'
            }

            post {
                // If Maven was able to run the tests, even if some of the test
                // failed, record the test results and archive the jar file.
                success {
                    echo 'Successfully Cloned Repository'
                }

                always {
                  echo "i tried..."
                }

                cleanup {
                  echo "after all other post condition"
                }
            }
        }

        // aws s3 에 파일을 올림
        stage('Deploy redteam') {
          steps {
            echo 'Deploy redteam project...'
            // redteam 프로젝트 s3에 올림, 이전에 반드시 EC2 instance profile 을 등록해야한다.
            //redteam 디렉토리의 모든 것들을  
            dir ('./redteam') {
                sh '''
                aws s3 sync ./ s3://sjabber
                '''
            }
          }
        }
    }
}
