package org.openindustries.ecosystem.station;

import io.quarkus.test.junit.QuarkusTest;
import org.junit.jupiter.api.Test;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.is;

@QuarkusTest
class MonitoringStationResourceTest {

    @Test
    void createAndTransition() {
        String id = given().contentType("application/json")
                .body("{\"name\":\"Flux Tower North\",\"stationType\":\"FLUX_TOWER\",\"readingsPerDay\":48}")
                .when().post("/api/v1/monitoring-stations")
                .then().statusCode(201)
                .body("status", is("INSTALLED"))
                .body("stationType", is("FLUX_TOWER"))
                .extract().path("id");

        given().when().post("/api/v1/monitoring-stations/" + id + "/activate")
                .then().statusCode(200).body("status", is("ACTIVE"));
    }
}
