package org.openindustries.ecosystem.station;

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

@Path("/api/v1/monitoring-stations")
@Produces(MediaType.APPLICATION_JSON)
@Tag(name = "Monitoring Stations", description = "Registry of monitoring stations feeding ecosystem condition data.")
public class MonitoringStationResource {

    @Inject
    MonitoringStationService service;

    @POST
    @Consumes(MediaType.APPLICATION_JSON)
    @Operation(summary = "Create a monitoringstation")
    public Response create(@Valid CreateMonitoringStationRequest request) {
        return Response.status(Response.Status.CREATED).entity(service.create(request)).build();
    }

    @GET
    @Operation(summary = "List all")
    public List<MonitoringStation> list() {
        return service.list();
    }

    @GET
    @Path("/{id}")
    @Operation(summary = "Get one by id")
    public MonitoringStation get(@PathParam("id") UUID id) {
        return service.get(id);
    }

    @POST
    @Path("/{id}/activate")
    @Operation(summary = "Activate (lifecycle transition, no body)")
    public MonitoringStation activate(@PathParam("id") UUID id) {
        return service.activate(id);
    }
}
