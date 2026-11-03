package org.openindustries.ecosystem.valuation;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import java.util.List;
import java.util.UUID;

@ApplicationScoped
public class ValuationRuleService {

    @Inject
    ValuationRuleRepository repository;

    public ValuationRule create(CreateValuationRuleRequest request) {
        return repository.save(ValuationRule.create(request.name(), request.serviceType(),
                request.maxValuePerHectareMinor(), request.maxUncertaintyPct(), request.priority()));
    }

    public ValuationRule get(UUID id) {
        return repository.findById(id)
                .orElseThrow(() -> new DomainException(404, "ValuationRule not found: " + id));
    }

    public List<ValuationRule> list() {
        return repository.findAll();
    }

    public ValuationRule activate(UUID id) {
        ValuationRule rule = get(id);
        rule.status = "ACTIVE";
        return repository.save(rule);
    }

    public Evaluation evaluate(EvaluateRequest request) {
        var candidates = repository.activeFor(request.serviceType());
        if (candidates.isEmpty()) {
            return new Evaluation("DENIED", null,
                    "No active rule for " + request.serviceType());
        }
        for (ValuationRule rule : candidates) {
            if (rule.maxUncertaintyPct <= request.requiredMaxUncertaintyPct()
                    && rule.maxValuePerHectareMinor >= request.requiredMaxValuePerHectareMinor()) {
                return new Evaluation("ALLOWED", rule.id,
                        "Rule '" + rule.name + "' covers the request");
            }
        }
        return new Evaluation("DEGRADED", candidates.get(0).id,
                "Active rules exist for " + request.serviceType() + " but none meets the requirement");
    }
}
