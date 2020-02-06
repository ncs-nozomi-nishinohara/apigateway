# apigateway

## nginx の設定ファイルを編集する機能

編集出来る内容

## 設定.yaml の内容

| Key                      | Overview                     | type   |
| :----------------------- | :--------------------------- | :----- |
| serverName               | サーバー名                   | string |
| locations                | location 設定(配列)          | list   |
| - url                    | proxy_pass url               | string |
| - location               | location URL                 | string |
| - prefix                 | プレフィックス               | string |
| proxy_set_header_default | デフォルト設定をするかどうか | bool   |
| - Host                   | デフォルトで設定される値     | string |
| - X-Real-IP              | デフォルトで設定される値     | string |
| - X-Forwarded-Proto      | デフォルトで設定される値     | string |
| - X-Forwarded-For        | デフォルトで設定される値     | string |
| proxy_set_header         | 各々で設定が可能             | map    |
| server_section           | 各々で設定が可能             | map    |

## Environment

| key               | Overview                             | Default                        |
| :---------------- | :----------------------------------- | :----------------------------- |
| NGINX_CONF        | 自動生成される内容を配置するフルパス | /etc/nginx/conf.d/default.conf |
| SETTING_FILE_NAME | 上記設定ファイルが配置されているパス | null                           |
