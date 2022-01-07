## API for Boilerplate

### Development
- Copy `.env.example` to `.env` and update according to requirement.
- Create `serviceAccountKey.json` file for firebase admin sdk.
- To run `docker-compose up` (with default configuration will run at `5000` and adminer runs at `5001`)

#### Run Boilerplate CLI 🖥  
- Run `docker-compose exec web sh` 
- After running type `./__debug_bin cli` you will start cli application. 
- Choose the commands to run afterwards.

#### Migration Commands 🛳  

| Command            | Desc                                           |
| -------------- | ---------------------------------------------- |
| `make migrate-up`   | runs migration up command                      |
| `make migrate-down` | runs migration down command                    |
| `make force`        | Set particular version but don't run migration |
| `make goto`         | Migrate to particular version                  |
| `make drop`         | Drop everything inside database                |
| `make create`       | Create new migration file(up & down)           |


**For Debugging 🐞** Debugger runs at `5002`. Vs code configuration is at `.vscode/launch.json` which will attach debugger to remote application.
