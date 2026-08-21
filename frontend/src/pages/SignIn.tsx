import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { authService } from '../services/authService';
import '../styles/auth.css';

export const SignIn = () => {
  const navigate = useNavigate();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await authService.signIn({ username, password });
      authService.saveToken(response.data.token);
      navigate('/main');
    } catch (error) {
      alert('Неверный логин или пароль');
    }
  };

  return (
    <div className="auth-container">
      <form onSubmit={handleSubmit} className="auth-form">
        <h2>Вход</h2>
        <input onChange={(e) => setUsername(e.target.value)} placeholder="Username" required />
        <div className="password-wrapper">
          <input 
            type={showPassword ? "text" : "password"} 
            onChange={(e) => setPassword(e.target.value)} 
            placeholder="Password" 
            required 
          />
          <button type="button" onClick={() => setShowPassword(!showPassword)}>
            {showPassword ? "Скрыть" : "Показать"}
          </button>
        </div>
        <button type="submit">Войти</button>
        <p className="auth-link">Нет аккаунта? <Link to="/sign-up">Зарегистрироваться</Link></p>
      </form>
    </div>
  );
};
