import { useState } from 'react';
import axios from 'axios';
import toast from 'react-hot-toast';

const SignIn = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await axios.post('http://localhost:8082/Sign-In', { username, password });
      toast.success('Login successful!');
    } catch (error) {
      console.error('Error:', error);
      toast.error('Login failed!');
    }
  };

  return (
    <form className="auth-container" onSubmit={handleSubmit}>
      <h2>Sign In</h2>
      <input type="text" placeholder="Username" onChange={(e) => setUsername(e.target.value)} />
      <input type="password" placeholder="Password" onChange={(e) => setPassword(e.target.value)} />
      <button type="submit">Sign In</button>
    </form>
  );
};

export default SignIn;
