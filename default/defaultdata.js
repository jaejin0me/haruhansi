const { MongoClient } = require('mongodb');

async function main() {
	const uri = "mongodb://localhost:27017?useUnifiedTopology=true";
	const client = new MongoClient(uri);

	try {
		// Connect to the MongoDB cluster
		await client.connect();

		// Make the appropriate DB calls
		await listDatabases(client);

	} catch (e) {
		console.error(e);
	} finally {
		await client.close();
	}
}

async function listDatabases(client) {
	databasesList = await client.db().admin().listDatabases();
	console.log("Databases:");
	databasesList.databases.forEach(db => console.log(` - ${db.name}`));
};

main().catch(console.error);