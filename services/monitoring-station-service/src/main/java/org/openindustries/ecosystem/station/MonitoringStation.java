package org.openindustries.ecosystem.station;

import java.util.UUID;

/** MonitoringStation aggregate — Registry of monitoring stations feeding ecosystem condition data. */
public class MonitoringStation {
    public UUID id;
        public String name;
        public String stationType;
        public int readingsPerDay;
    public String status;

    public static MonitoringStation create(String name, String stationType, int readingsPerDay) {
        MonitoringStation e = new MonitoringStation();
        e.id = UUID.randomUUID();
        e.name = name;
        e.stationType = stationType;
        e.readingsPerDay = readingsPerDay;
        e.status = "INSTALLED";
        return e;
    }
}
