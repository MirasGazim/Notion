import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { authService } from '../services/authService';
import '../styles/auth.css';

export const SignUp = () => {
  const navigate = useNavigate();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await authService.signUp({ email, username, password });
      authService.saveToken(response.data.token);
      navigate('/main');
    } catch (error) {
      alert('Ошибка при регистрации');
    }
  };

  return (
    <div className="auth-container">
      <form onSubmit={handleSubmit} className="auth-form">
        <h2>Регистрация</h2>
        <input type="email" onChange={(e) => setEmail(e.target.value)} placeholder="Email" required />
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
        <button type="submit">Зарегистрироваться</button>
        <p className="auth-link">Уже есть аккаунт? <Link to="/sign-in">Войти</Link></p>
      </form>
    </div>
  );
};
