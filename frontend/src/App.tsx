import {
  Navigate,
  Route,
  Routes,
} from 'react-router-dom';

import { CompanyRegisterPage } from './pages/Register/CompanyRegisterPage';
import { UserRegisterPage } from './pages/Register/UserRegisterPage';
import { VirtualOfficePage } from './pages/VirtualOffice/VirtualOfficePage';

function App() {
  return (
    <Routes>
      <Route
        path="/"
        element={<Navigate to="/register/company" replace />}
      />

      <Route
        path="/register/company"
        element={<CompanyRegisterPage />}
      />

      <Route
        path="/register/company/:companyId/user"
        element={<UserRegisterPage />}
      />

      <Route
        path="/office"
        element={<VirtualOfficePage />}
      />
    </Routes>
  );
}

export default App;