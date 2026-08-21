import { Sidebar } from '../sidebar/Sidebar';
import '../../styles/layout/main.css';

export const MainLayout = ({ children, onWorkspaceSelect }: { children: React.ReactNode, onWorkspaceSelect: (id: string) => void }) => {
  return (
    <div className="app-container">
      <Sidebar onWorkspaceSelect={onWorkspaceSelect} />
      <div className="main-content">
        <div className="content-wrapper">
          {children}
        </div>
      </div>
    </div>
  );
};
