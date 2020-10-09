# Recomator

**This is not an officially supported Google product.**

Google Cloud Platform shows recommendations for underutilized/overutilized resources. Recomator extends these functionalities by
allowing to also apply recommendations with one click and by giving an overview of recommendations for all projects. 
Customers will be able to deploy a premade package on one of the serverless solutions and run it in AppEngine in their own GCP account.

<!-- This table of contents is auto-generated using Markdown All in One (VS Code extension). If you have it installed, it updates on saves. -->
- [Recomator](#recomator)
  - [Frontend config](#frontend-config)
  - [Frontend commands](#frontend-commands)
    - [Install dependencies](#install-dependencies)
    - [Compile and hot-reload for development](#compile-and-hot-reload-for-development)
    - [Compile and minify for production](#compile-and-minify-for-production)
    - [Run your unit tests](#run-your-unit-tests)
    - [Lint and fix files](#lint-and-fix-files)
    - [Customize Vue configuration](#customize-vue-configuration)
  - [How to deploy](#how-to-deploy)
    - [Create your application on App Engine](#create-your-application-on-app-engine)
    - [Install Google Cloud SDK](#install-google-cloud-sdk)
    - [Create the directory](#create-the-directory)
    - [Create credentials](#create-credentials)
    - [Deploy the app](#deploy-the-app)
  - [Source Code Headers](#source-code-headers)

## Frontend config

`src/config.ts` contains two web addresses that you should provide:
For testing purposes (using `npm run serve`), use the first one. 
For production (`npm run build`), use the second one. 
```typescript
  public static DEVELOPMENT_BACKEND_ADDRESS = "http://dev.yourbackend.com";
  public static PRODUCTION_BACKEND_ADDRESS = "http://yourbackend.com";
```

## Frontend commands
Note: The following commands need to be run in the `./frontend/` directory:

### Install dependencies

```
npm install
```

### Compile and hot-reload for development
```
npm run serve
```

### Compile and minify for production
```
npm run build
```

### Run your unit tests
```
npm run test:unit
```

### Lint and fix files
```
npm run lint
```

### Customize Vue configuration
See [Configuration Reference](https://cli.vuejs.org/config/).

## How to deploy

### Create your application on App Engine

- Create a [project on GCP](https://console.cloud.google.com/projectcreate).

- Make sure that the APIs listed below are **Enabled**:
  - [Cloud Resource Manager API](https://console.cloud.google.com/apis/library/cloudresourcemanager.googleapis.com), 
  - [Recommender API](https://console.cloud.google.com/apis/library/recommender.googleapis.com),
  - [Service Usage API](https://pantheon.corp.google.com/apis/library/serviceusage.googleapis.com), 
  - [Compute Engine API](https://pantheon.corp.google.com/apis/library/compute.googleapis.com).
  
- Go to [App Engine](http://console.cloud.google.com/appengine) page.

- Click `create application`.

![Image-Create](https://github.com/recomator/docs/create-application.png)

- Choose preferred region, language - `Go`, environment - `standard`.

### Install Google Cloud SDK

Note that you don’t need to install it, if you are using [Cloud Shell](https://cloud.google.com/shell).
Otherwise, to install you need to follow this tutorial https://cloud.google.com/sdk/docs/install

### Create the directory

Run the following commands:
```
git clone https://github.com/googleinterns/recomator
```
```
cd recomator
```
```
gcloud init 
```
You will be asked to choose an account and project. Choose the project in which you created the app.

Now run 
```
gcloud app describe | grep defaultHostname
```
You will see the value of defaultHosthame (your future app address).

Copy it, you’ll need it in the next step.

### Create credentials
 - Go to APIs & Services -> OAuth consent screen.
 
 ![Image-Consent-Screen](https://github.com/recomator/docs/consent-screen.png)
 
 - Choose your app’s name, support email.
 
 - Click `Add domain` and add the address of your app (from previous section, last step).
 
 - Add developer’s email address and click `Save & continue`.
 
 - Add scopes: 
```
.../auth/userinfo.email
.../auth/cloud-platform
```

- Click `Save`, then `Back to dashboard`.

- Go to Credentials-> Create Credentials -> OAuth client ID

![Image-OAuth-Client](https://github.com/recomator/docs/oauth-client.png)

- Choose `Web application`

- Add authorized URI: `https://<YOUR APP ADDRESS>/auth`, where `<YOUR APP ADDRESS>` is the address of your application.

- Copy clientID and clientSecret.

### Deploy the app

- Install npm https://www.npmjs.com/get-npm
- Go to your app directory (the one where you cloned the recomator repository).
- Set up some environmental variables:
```
export CLIENT_ID=<YOUR CLIENT ID>
```
```
export CLIENT_SECRET=<YOUR CLIENT SECRET>
```
```
export APP_ADDRESS=<YOUR ADDRESS>
```

- Set up configs:
```
echo "{\"clientID\":\"$CLIENT_ID\",\"clientSecret\":\"$CLIENT_SECRET\",\"redirectURL\":\"https://$APP_ADDRESS/auth\"}" > config.json
```
```
echo "VUE_APP_BACKEND_ADDRESS=https://$APP_ADDRESS/api" > frontend/.env
```

- Build frontend:

```
cd frontend
```
```
npm install
```
```
npm run build
```
 
- Deploy the app:
```
cd ..
```
```
cat deploy_gcloudignore.txt >> .gcloudignore
```
```
gcloud app deploy
```


## Source Code Headers

Every file containing source code must include copyright and license
information. This includes any JS/CSS files that you might be serving out to
browsers. (This is to help well-intentioned people avoid accidental copying that
doesn't comply with the license.)

Apache header:

    Copyright 2020 Google LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        https://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

