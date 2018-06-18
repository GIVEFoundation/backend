TOKEN="4b6dbabba1ee5f2e5f3591201373164b"
curl -s -H "Content-Type: application/json" -H "GIVEAPIToken: $TOKEN" -X POST -d '{"id":"0x9AbAB02EcBe8A917C266681B37d1f45f56191bDb", "name":"Danny Wood", "date_of_birth":"2010-11-10", "parents_emails":["pa@gmail.com","ma@gmail.com"], "school_name":"1st School of Hawaii", "id_tag_name":"Big Danny"}' 'http://localhost:8080/api/v1/give/kids'
