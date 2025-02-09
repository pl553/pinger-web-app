import { useEffect, useState } from "react"
import { getContainersList } from "../api/containers-list";

type Container = {
  ContainerIp: string,
  PingTimeMs: number,
  LastSuccessfulPing: string | null,
  Status: string,
}

export function ContainersListPage() {
  const [containers, setContainers] = useState<Container[]>([]);

  const fetchContainers = async () => {
    try {
      const data = await getContainersList()
      data.sort((a, b) => a.ContainerIp > b.ContainerIp ? 1 : -1)
      setContainers(data)
    } catch (error) {
      console.error("Error fetching containers:", error);
    }
  }

  useEffect(() => {
    fetchContainers();

    const intervalId = setInterval(fetchContainers, 1000);

    return () => clearInterval(intervalId);
  }, []);

  const formatDate = (ts: string) => {
    const date = new Date(ts);
    return new Intl.DateTimeFormat("en-GB", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false,
    }).format(date);
  };

  return (
    <div className="container mx-auto py-2">
      <table className="min-w-full border-gray-200">
        <thead className="bg-gray-100">
          <tr>
            <th className="border">Container IP</th>
            <th className="border">Ping Time (ms)</th>
            <th className="border">Last Successful Ping Date</th>
            <th className="py-2 border">Status</th>
          </tr>
        </thead>
        <tbody>
          {containers.map((container) => (
            <tr key={container.ContainerIp} className="hover:bg-gray-50">
              <td className="py-2 px-4 border">{container.ContainerIp}</td>
              <td className="py-2 px-4 border">{container.PingTimeMs}</td>
              <td className="py-2 px-4 border">{container.LastSuccessfulPing && formatDate(container.LastSuccessfulPing)}</td>
              <td
                className={`border px-4 py-2 font-semibold ${
                    container.Status === "UP" ? "background-color text-green-600" : "text-red-600"
                }`}
              >
                {container.Status}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}