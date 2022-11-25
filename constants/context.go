package constants

const (
	// DBTransaction is database transaction handle set at router context
	DBTransaction = "db_trx"

	// Claims -> authentication claims
	Claims = "Claims"

	// UID -> authenticated user's id
	UID = "UID"

	//DUMMYADMIN ->
	DUMMYADMIN = "Administrator"

	//DUMMYEMAIL ->
	DUMMYEMAIL = "dummyrtw@mailinator.com"

	FirebaseUseLoginUrl = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key="
	FirebaseUID         = "firebase_uid"
)
