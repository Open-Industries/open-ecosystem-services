package org.openindustries.ecosystem.functions.health;

import io.quarkus.test.junit.QuarkusTest;
import org.junit.jupiter.api.Test;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.is;

@QuarkusTest
class FunctionTest {

    @Test
    void healthyIsNormal() {
        given().contentType("application/json")
                .body("{\"assetId\":\"a-1\",\"deforestationIndex\":3,\"pollutionIndex\":5,\"invasiveIndex\":2}")
                .when().post("/ecosystemHealthScorer")
                .then().statusCode(200).body("decision", is("NORMAL")).body("score", is(0));
    }

    @Test
    void stackedSignalsAreCritical() {
        given().contentType("application/json")
                .body("{\"assetId\":\"a-2\",\"deforestationIndex\":85,\"pollutionIndex\":70,\"invasiveIndex\":65}")
                .when().post("/ecosystemHealthScorer")
                .then().statusCode(200).body("decision", is("CRITICAL")).body("score", is(100));
    }
}
