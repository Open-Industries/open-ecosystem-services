package org.openindustries.ecosystem.functions.delivery;

import io.quarkus.test.junit.QuarkusTest;
import org.junit.jupiter.api.Test;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.is;

@QuarkusTest
class FunctionTest {

    @Test
    void compliant() {
        given().contentType("application/json")
                .body("{\"assetId\":\"a-1\",\"measuredDegradationPct\":2,\"targetDegradationPct\":5,\"measuredServiceIndexPct\":91.0,\"targetServiceIndexPct\":85.0}")
                .when().post("/serviceDeliveryEvaluator")
                .then().statusCode(200).body("compliant", is(true)).body("findings.size()", is(0));
    }

    @Test
    void breachIsReported() {
        given().contentType("application/json")
                .body("{\"assetId\":\"a-2\",\"measuredDegradationPct\":14,\"targetDegradationPct\":5,\"measuredServiceIndexPct\":91.0,\"targetServiceIndexPct\":85.0}")
                .when().post("/serviceDeliveryEvaluator")
                .then().statusCode(200).body("compliant", is(false)).body("findings.size()", is(1));
    }
}
