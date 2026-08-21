import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
import SignIn from './SignIn';
import SignUp from './SignUp';

function App() {
  return (
    <Router>
      <Toaster position="top-right" />
      <div className="layout">
        <div className="sidebar">
          <div className="logo">
            <div className="logo-icon">N</div>
            Notion Clone
          </div>
          <div className="sidebar-section">
            <Link to="/sign-in" className="nav-link">Sign In</Link>
            <Link to="/sign-up" className="nav-link">Sign Up</Link>
          </div>
          <div className="sidebar-section">
            <h3 style={{ fontSize: '0.7rem', color: '#999', paddingLeft: '12px' }}>Workspace</h3>
            <div style={{ padding: '8px 12px', fontSize: '0.85rem', color: '#999' }}>No pages yet...</div>
          </div>
          <div className="footer">
            Built with React & Go <br />
            v1.0.0
          </div>
        </div>
        <div className="main-content">
          <Routes>
            <Route path="/sign-in" element={<SignIn />} />
            <Route path="/sign-up" element={<SignUp />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
}

export default App;
