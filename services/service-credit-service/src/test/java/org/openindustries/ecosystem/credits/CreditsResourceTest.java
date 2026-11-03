package org.openindustries.ecosystem.credits;

import io.quarkus.test.junit.QuarkusTest;
import org.junit.jupiter.api.Test;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.is;

@QuarkusTest
class CreditsResourceTest {

    @Test
    void consumptionDrawsDown() {
        String id = given().contentType("application/json")
                .body("{\"assetId\":\"11111111-1111-1111-1111-111111111111\",\"capacityUnits\":10000}")
                .when().post("/api/v1/credit-accounts")
                .then().statusCode(201).body("capacityUnits", is(10000))
                .extract().path("id");

        given().contentType("application/json")
                .body("{\"serviceType\":\"CARBON\",\"units\":1500}")
                .when().post("/api/v1/credit-accounts/" + id + "/issuances")
                .then().statusCode(200)
                .body("remainingAfter", is(8500));
    }

    @Test
    void overdrawIsRejected() {
        String id = given().contentType("application/json")
                .body("{\"assetId\":\"11111111-1111-1111-1111-111111111111\",\"capacityUnits\":10000}")
                .when().post("/api/v1/credit-accounts")
                .then().statusCode(201)
                .extract().path("id");

        given().contentType("application/json")
                .body("{\"serviceType\":\"CARBON\",\"units\":50000}")
                .when().post("/api/v1/credit-accounts/" + id + "/issuances")
                .then().statusCode(422).body("status", is(422));

        given().when().get("/api/v1/credit-accounts/" + id)
                .then().statusCode(200).body("capacityUnits", is(10000));
    }
}
