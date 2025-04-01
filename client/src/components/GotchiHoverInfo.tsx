import { useMemo } from "react";
import "./GotchiHoverInfo.css";

interface GotchiData {
  name?: string;
  gotchiId?: string;
  stats?: {
    attack?: number;
    harvest_duration_s?: number;
    hp_current?: number;
    hp_max?: number;
    trade_duration_s?: number;
  };
  inventory?: { [key: string]: number };
  personality?: string[];
  activityLog?: Array<{
    Description?: string;
    LogTime?: string;
  }>;
}

interface GotchiHoverInfoProps {
  data: GotchiData | null;
}

export const GotchiHoverInfo = ({ data }: GotchiHoverInfoProps) => {
  if (!data) {
    return (
      <div className="gotchi-hover-info">
        <h3 className="gotchi-section-title">No Data Available</h3>
        <p>Entity data could not be loaded.</p>
      </div>
    );
  }

  const formattedLog = useMemo(() => {
    if (!data.activityLog) return [];
    return data.activityLog.map((entry) => ({
      description: entry?.Description ?? "Unknown activity",
      time: entry?.LogTime ? new Date(entry.LogTime).toLocaleString() : "Unknown time",
    }));
  }, [data.activityLog]);

  const personalityColors: { [key: string]: string } = {
    Turnt: "#ff4444",
    Based: "#ff8844",
    Cuddly: "#44ff88",
    Mindbender: "#4488ff",
    Soulful: "#ff44ff",
  };

  return (
    <div className="gotchi-hover-info">
      {/* Avatar Section */}
      <div className="gotchi-avatar-section">
        <img
          src={
            data.gotchiId
              ? `data:image/svg+xml;base64,${btoa(
                  `<svg width="64" height="64"><rect width="64" height="64" fill="#ccc"/></svg>`
                )}`
              : `data:image/svg+xml;base64,${btoa(
                  `<svg width="64" height="64"><rect width="64" height="64" fill="#999"/></svg>`
                )}`
          }
          alt={data.name ?? "Unknown"}
          className="gotchi-avatar"
        />
        <h2 className="gotchi-name">{data.name ?? "Unknown Entity"}</h2>
      </div>

      {/* Stats Section */}
      <div className="gotchi-section">
        <h3 className="gotchi-section-title">Stats</h3>
        <div className="gotchi-stats">
          <div className="gotchi-stat">
            <span className="stat-label">HP:</span>
            <span className="stat-value">
              {data.stats?.hp_current ?? 0}/{data.stats?.hp_max ?? 0}
            </span>
          </div>
          <div className="gotchi-stat">
            <span className="stat-label">Attack:</span>
            <span className="stat-value">{data.stats?.attack ?? 0}</span>
          </div>
          <div className="gotchi-stat">
            <span className="stat-label">Harvest Time:</span>
            <span className="stat-value">{data.stats?.harvest_duration_s ?? 0}s</span>
          </div>
          <div className="gotchi-stat">
            <span className="stat-label">Trade Time:</span>
            <span className="stat-value">{data.stats?.trade_duration_s ?? 0}s</span>
          </div>
        </div>
      </div>

      {/* Personality Section */}
      <div className="gotchi-section">
        <h3 className="gotchi-section-title">Personality</h3>
        <div className="gotchi-personality">
            {/* Explicit check for undefined or empty personality */}
            {data.personality && data.personality.length > 0 ? (
            data.personality.map((trait, index) => (
                <span
                key={index}
                className="personality-trait"
                style={{ backgroundColor: personalityColors[trait] ?? "#888" }}
                >
                {trait}
                </span>
            ))
            ) : (
            <span className="personality-trait" style={{ backgroundColor: "#888" }}>
                None
            </span>
            )}
        </div>
    </div>

      {/* Inventory Section */}
      <div className="gotchi-section gotchi-inventory-section">
        <h3 className="gotchi-section-title">Inventory</h3>
        <div className="gotchi-inventory">
          {data.inventory && Object.keys(data.inventory).length > 0 ? (
            Object.entries(data.inventory).map(([item, count], index) => (
              <div key={index} className="inventory-item">
                <img
                  src={`data:image/svg+xml;base64,${btoa(
                    `<svg width="32" height="32"><rect width="32" height="32" fill="${
                      item === "fomoberry" ? "#ff4444" : "#888"
                    }"/></svg>`
                  )}`}
                  alt={item}
                  className="inventory-icon"
                />
                <span className="inventory-count">{count}</span>
                <div className="inventory-tooltip">
                  <h4>{item.charAt(0).toUpperCase() + item.slice(1)}</h4>
                  <p>Basic resource item</p>
                </div>
              </div>
            ))
          ) : (
            <p>No items in inventory</p>
          )}
        </div>
      </div>

      {/* Activity Log Section */}
      <div className="gotchi-section gotchi-activity-section">
        <h3 className="gotchi-section-title">Activity Log</h3>
        <div className="gotchi-activity">
          {formattedLog.length > 0 ? (
            formattedLog.map((entry, index) => (
              <div key={index} className="activity-entry">
                <span className="activity-time">{entry.time}</span>
                <span className="activity-description">{entry.description}</span>
              </div>
            ))
          ) : (
            <p>No activities recorded</p>
          )}
        </div>
      </div>
    </div>
  );
};
