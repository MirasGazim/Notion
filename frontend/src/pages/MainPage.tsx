import { useState } from 'react';
import { MainLayout } from '../components/layout/MainLayout';
import { authService } from '../services/authService';

export const MainPage = () => {
  const [selectedWorkspaceId, setSelectedWorkspaceId] = useState<string | null>(null);
  const [blocks, setBlocks] = useState([]);

  const handleWorkspaceSelect = async (id: string) => {
    setSelectedWorkspaceId(id);
    try {
      const response = await authService.getWorkspaceBlocks(id);
      // Теперь правильно получаем массив блоков
      setBlocks(response.data.blocks || []);
    } catch (error) {
      console.error('Ошибка загрузки блоков:', error);
      setBlocks([]); // Очищаем при ошибке
    }
  };

  return (
    <MainLayout onWorkspaceSelect={handleWorkspaceSelect}>
      {selectedWorkspaceId ? (
        <div>
          {blocks.length === 0 ? (
            <div style={{ color: '#767671', marginTop: '20px' }}>
              <p>Здесь пока пусто. Создайте первый блок, чтобы начать работу.</p>
            </div>
          ) : (
            <div>
              <h1>Содержимое пространства: {selectedWorkspaceId}</h1>
              <pre>{JSON.stringify(blocks, null, 2)}</pre>
            </div>
          )}
        </div>
      ) : (
        <h1>Выберите рабочее пространство из боковой панели</h1>
      )}
    </MainLayout>
  );
};
