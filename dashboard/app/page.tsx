"use client";

import { useCallback, useEffect, useState } from "react";
import DockerTab from "./components/DockerTab";
import EnvTab from "./components/EnvTab";
import NginxTab from "./components/NginxTab";
import OverviewTab from "./components/OverviewTab";
import PortsTab from "./components/PortsTab";
import Sidebar, { type TabId } from "./components/Sidebar";
import SSHTab from "./components/SSHTab";
import ThemeLangToggle from "./components/ThemeLangToggle";
import { dockerApi, envApi, portsApi, sshApi } from "./lib/api";
import type { DockerContainer, EnvVariable, PortEntry, SSHHost } from "./types";

export default function Home() {
  const [activeTab, setActiveTab] = useState<TabId>("overview");
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  const [containers, setContainers] = useState<DockerContainer[]>([]);
  const [ports, setPorts] = useState<PortEntry[]>([]);
  const [sshHosts, setSSHHosts] = useState<SSHHost[]>([]);
  const [envVars, setEnvVars] = useState<EnvVariable[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const [portsData, dockerData, sshData, envData] =
        await Promise.allSettled([
          portsApi.scan(),
          dockerApi.list(),
          sshApi.list(),
          envApi.show(),
        ]);

      if (portsData.status === "fulfilled") setPorts(portsData.value);
      if (dockerData.status === "fulfilled") setContainers(dockerData.value);
      if (sshData.status === "fulfilled") setSSHHosts(sshData.value);
      if (envData.status === "fulfilled") {
        const vars = Object.entries(envData.value).map(([key, value]) => ({
          key,
          value: value as string,
          isSensitive: /secret|password|token|key|api/i.test(key),
        }));
        setEnvVars(vars);
      }
    } catch (err) {
      console.error("Failed to fetch data:", err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <div className="min-h-screen bg-[var(--color-background)]">
      <Sidebar
        activeTab={activeTab}
        onTabChange={setActiveTab}
        collapsed={sidebarCollapsed}
        onToggle={() => setSidebarCollapsed(!sidebarCollapsed)}
      />

      <main
        className={`transition-all duration-300 min-h-screen ${
          sidebarCollapsed ? "ml-[68px]" : "ml-[240px]"
        }`}
      >
        <div className="flex justify-end p-4 border-b border-[var(--color-border)]">
          <ThemeLangToggle />
        </div>
        <div className="p-6 max-w-[1400px]">
          {loading ? (
            <div className="flex items-center justify-center h-[60vh]">
              <div className="text-[var(--color-muted)] text-sm">
                ⏳ Đang tải dữ liệu từ CLI...
              </div>
            </div>
          ) : (
            <>
              {activeTab === "overview" && (
                <OverviewTab
                  containers={containers}
                  ports={ports}
                  sshHosts={sshHosts}
                  envVarCount={envVars.length}
                  onNavigate={setActiveTab}
                />
              )}
              {activeTab === "docker" && <DockerTab containers={containers} />}
              {activeTab === "ports" && <PortsTab ports={ports} />}
              {activeTab === "ssh" && <SSHTab hosts={sshHosts} />}
              {activeTab === "env" && <EnvTab envVars={envVars} />}
              {activeTab === "nginx" && <NginxTab configs={[]} />}
            </>
          )}
        </div>
      </main>
    </div>
  );
}
