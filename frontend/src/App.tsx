import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { SignIn } from './pages/SignIn';
import { SignUp } from './pages/SignUp';
import { MainPage } from './pages/MainPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/sign-in" element={<SignIn />} />
        <Route path="/sign-up" element={<SignUp />} />
        <Route path="/main" element={<MainPage />} />
        <Route path="/" element={<Navigate to="/sign-in" />} />
      </Routes>
    </Router>
  );
}

export default App;
