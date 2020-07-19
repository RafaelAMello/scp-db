mv scp.db localdb/

docker run --rm -p 3000:3000 \
-v $(pwd)/metabase-data:/metabase-data \
-v $(pwd)/localdb:/localdb \
-e "MB_DB_FILE=/metabase-data/metabase.db" \
--name metabase metabase/metabase

docker run -d -p 5432:5432 \
    -e POSTGRES_USER=scp \
    -e POSTGRES_PASSWORD=shh \
    --name postgres postgres:alpine

docker run \
    --name neo4j \
    --rm \
    -p 7474:7474 -p 7687:7687 \
    -v $(pwd)/neo4j/data:/data \
    -v $(pwd)/neo4j/logs:/logs \
    -v $(pwd)/neo4j/import:/var/lib/neo4j/import \
    -v $(pwd)/neo4j/plugins:/plugins \
    --env NEO4J_AUTH=neo4j/password \
    neo4j:latest

export GOOGLE_APPLICATION_CREDENTIALS="$(pwd)/scp-db.json"