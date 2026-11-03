package org.openindustries.ecosystem.credits;

import jakarta.inject.Inject;
import jakarta.validation.Valid;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import org.eclipse.microprofile.openapi.annotations.Operation;
import org.eclipse.microprofile.openapi.annotations.tags.Tag;

import java.util.List;
import java.util.UUID;

@Path("/api/v1/credit-accounts")
@Produces(MediaType.APPLICATION_JSON)
@Tag(name = "Service Credits", description = "Verified service-capacity accounts per asset and the credit issuances against them.")
public class CreditsResource {

    @Inject
    CreditsService service;

    @POST
    @Consumes(MediaType.APPLICATION_JSON)
    @Operation(summary = "Open a creditaccount")
    public Response open(@Valid OpenCreditAccountRequest request) {
        return Response.status(Response.Status.CREATED).entity(service.open(request)).build();
    }

    @GET
    @Path("/{id}")
    @Operation(summary = "Get one by id")
    public CreditAccount get(@PathParam("id") UUID id) {
        return service.get(id);
    }

    @POST
    @Path("/{id}/issuances")
    @Consumes(MediaType.APPLICATION_JSON)
    @Operation(summary = "Draw from the account (rejects if it would overdraw, 422)")
    public CreditIssuance consume(@PathParam("id") UUID id, @Valid CreateCreditIssuanceRequest request) {
        return service.consume(id, request);
    }

    @GET
    @Path("/{id}/issuances")
    @Operation(summary = "List consumption records of an account")
    public List<CreditIssuance> records(@PathParam("id") UUID id) {
        return service.recordsOf(id);
    }
}
