package org.openindustries.ecosystem.asset;

import io.quarkus.test.junit.QuarkusTest;
import org.junit.jupiter.api.Test;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.endsWith;
import static org.hamcrest.CoreMatchers.is;

@QuarkusTest
class EcosystemAssetResourceTest {

    @Test
    void createAndFetch() {
        String id = given().contentType("application/json")
                .body("{\"name\":\"Laguna Verde\",\"assetType\":\"WETLAND\",\"hectares\":320,\"regionCode\":\"MX-YUC-02\"}")
                .when().post("/api/v1/ecosystem-assets")
                .then().statusCode(201)
                .body("status", is("ACTIVE"))
                .body("assetType", is("WETLAND"))
                .extract().path("id");

        given().when().get("/api/v1/ecosystem-assets/" + id)
                .then().statusCode(200).body("status", is("ACTIVE"));
    }

    @Test
    void transitionWithoutBody() {
        // Lifecycle transitions take no body and must not require a Content-Type (no 415).
        String id = given().contentType("application/json")
                .body("{\"name\":\"Laguna Verde\",\"assetType\":\"WETLAND\",\"hectares\":320,\"regionCode\":\"MX-YUC-02\"}")
                .when().post("/api/v1/ecosystem-assets")
                .then().statusCode(201)
                .extract().path("id");

        given().when().post("/api/v1/ecosystem-assets/" + id + "/protect")
                .then().statusCode(200).body("status", is("PROTECTED"));
    }
}
