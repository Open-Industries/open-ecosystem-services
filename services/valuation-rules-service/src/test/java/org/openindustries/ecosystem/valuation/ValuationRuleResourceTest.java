package org.openindustries.ecosystem.valuation;

import io.quarkus.test.junit.QuarkusTest;
import org.junit.jupiter.api.Test;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.is;

@QuarkusTest
class ValuationRuleResourceTest {

    @Test
    void activateThenEvaluateAllowed() {
        String id = given().contentType("application/json")
                .body("{\"name\":\"carbon-forest\",\"serviceType\":\"CARBON\",\"maxValuePerHectareMinor\":5000,\"maxUncertaintyPct\":15,\"priority\":1}")
                .when().post("/api/v1/valuation-rules")
                .then().statusCode(201).body("status", is("DRAFT"))
                .extract().path("id");

        given().when().post("/api/v1/valuation-rules/" + id + "/activate")
                .then().statusCode(200).body("status", is("ACTIVE"));

        given().contentType("application/json")
                .body("{\"serviceType\":\"CARBON\",\"requiredMaxUncertaintyPct\":20,\"requiredMaxValuePerHectareMinor\":4000}")
                .when().post("/api/v1/valuation-rules/evaluate")
                .then().statusCode(200).body("decision", is("ALLOWED"));
    }

    @Test
    void evaluateWithoutRuleIsDenied() {
        given().contentType("application/json")
                .body("{\"serviceType\":\"WATER\",\"requiredMaxUncertaintyPct\":20,\"requiredMaxValuePerHectareMinor\":4000}")
                .when().post("/api/v1/valuation-rules/evaluate")
                .then().statusCode(200).body("decision", is("DENIED"));
    }
}
