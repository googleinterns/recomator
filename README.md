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

