## DB Setup
### Step 1:
In the root folder run the command
```bash
podman-compose -f doc/docker/compose.yml up -d
```

### Step 2:
Login to the database using root access, set the host machine port 3386

### Step 3
Run the queries described in the queries.db file in a MySQL console
