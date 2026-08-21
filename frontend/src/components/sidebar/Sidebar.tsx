import { useEffect, useState } from 'react';
import { CreateWorkspace } from './CreateWorkspace';
import { authService } from '../../services/authService';
import '../../styles/layout/main.css';

export const Sidebar = ({ onWorkspaceSelect }: { onWorkspaceSelect: (id: string) => void }) => {
  const [workspaces, setWorkspaces] = useState([]);

  const fetchWorkspaces = async () => {
    try {
      const response = await authService.getAllWorkspaces();
      setWorkspaces(response.data);
    } catch (error) {
      console.error('Ошибка загрузки пространств:', error);
    }
  };

  useEffect(() => {
    fetchWorkspaces();
  }, []);

  return (
    <div className="sidebar">
      <div className="sidebar-header">Workspaces</div>
      {workspaces.map((ws: any) => (
        <div key={ws.id} className="sidebar-item" onClick={() => onWorkspaceSelect(ws.id)}>
          {ws.name}
        </div>
      ))}
      <CreateWorkspace onCreated={fetchWorkspaces} />
    </div>
  );
};
