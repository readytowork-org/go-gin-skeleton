## Boilerplate API

Boilerplate API template includes all the common packages and setup used for API development in this Company.

### Development

- Copy `.env.example` to `.env` and update according to requirement.
- Create `serviceAccountKey.json` file for firebase admin sdk and copy serviceAccountKey.json.example file content for basic setup.
- To run `docker-compose up` (with default configuration will run at `5000` and adminer runs at `5001`)

#### Run Boilerplate CLI 🖥

- Run `docker-compose exec web sh`
- After running type `./__debug_bin cli` you will start cli application.
- Choose the commands to run afterwards.

#### Migration Commands 🛳

| Command                       | Desc                                           |
| ----------------------------  | ---------------------------------------------- |
| `make migrate-status`         | See status for migration file status           |
| `make create`                 | Create new migration file(up & down)           |
| `make migrate-up`             | runs migration up command                      |
| `make migrate-down`           | runs migration down command                    |
| `make force`                  | Set particular version but don't run migration |
| `make goto`                   | Migrate to particular version                  |
| `make drop`                   | Drop everything inside database                |
| `make crud`                   | Create CRUD template                           |

### Implemented Feature

- Dependency Injection (go-fx)
- Routing (gin web framework)
- Environment Files
- Logging (file saving on production) zap
- Middlewares (cors)
- Middleware transaction
- Database Setup (mysql)
- Models Setup and Automigrate (gorm)
- Repositories
- Implementing Basic CRUD Operation
- CRUD Scaffold Generator
- Migration Runner Implementation
- Live code refresh
- Firebase basic setup.
- GCP Cloud Storage bucket setup.
- Twilio Basic setup.
- Gmail Api Setup.

**For Debugging 🐞** Debugger runs at `5002`. Vs code configuration is at `.vscode/launch.json` which will attach debugger to remote application.

## For auto generate of CRUD(Create, ReaD, Update & Delete) api following informations are needed and will be asked in terminal:


- resource-name: name of CRUD in upper camelCase. examples:Food,Puppy,ProductCategory etc.  

- resource-table-name: name of CRUD in lower snake case. examples:food,puppy,product_category etc.

- plural-resource-table-name: plural name for the table going to be created. example: foods, puppies, product_categories.

- plural-resource-name: plural name of CRUD in Upper camelCase. examples:Foods,Puppies,ProductCategories etc.
