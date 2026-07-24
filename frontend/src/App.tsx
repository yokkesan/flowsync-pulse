import {
  Navigate,
  Route,
  Routes,
} from 'react-router-dom';

import { ProtectedRoute } from './components/auth/ProtectedRoute';
import { LoginPage } from './pages/Login/LoginPage';
import { CompanyRegisterPage } from './pages/Register/CompanyRegisterPage';
import { UserRegisterPage } from './pages/Register/UserRegisterPage';
import { VirtualOfficePage } from './pages/VirtualOffice/VirtualOfficePage';
import { ProjectListPage } from './pages/Projects/ProjectListPage';
import { ProjectCreatePage } from './pages/Projects/ProjectCreatePage';
import { ProjectDetailPage } from './pages/Projects/ProjectDetailPage';
import { ProjectEditPage } from './pages/Projects/ProjectEditPage';

import { AllTaskListPage } from './pages/Tasks/AllTaskListPage';
import { TaskListPage } from './pages/Tasks/TaskListPage';
import { TaskCreatePage } from './pages/Tasks/TaskCreatePage';
import { TaskDetailPage } from './pages/Tasks/TaskDetailPage';
import { TaskEditPage } from './pages/Tasks/TaskEditPage';

function App() {
  return (
    <Routes>
      <Route
        path="/"
        element={
          <Navigate
            to="/register/company"
            replace
          />
        }
      />

      <Route
        path="/login"
        element={<LoginPage />}
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
        element={
          <ProtectedRoute>
            <VirtualOfficePage />
          </ProtectedRoute>
        }
      />

      <Route
        path="/projects"
        element={
          <ProtectedRoute>
            <ProjectListPage />
          </ProtectedRoute>
        }
      />

      <Route
        path="/projects/new"
        element={
          <ProtectedRoute>
            <ProjectCreatePage />
          </ProtectedRoute>
        }
      />

      <Route
        path="/projects/:projectId"
        element={
          <ProtectedRoute>
            <ProjectDetailPage />
          </ProtectedRoute>
        }
      />

      <Route
        path="/projects/:projectId/edit"
        element={
          <ProtectedRoute>
            <ProjectEditPage />
          </ProtectedRoute>
        }
      />

      <Route
        path="/projects/:projectId/tasks"
        element={
          <ProtectedRoute>
            <TaskListPage />
          </ProtectedRoute>
        }
      />

      <Route
        path="/tasks"
        element={
          <ProtectedRoute>
            <AllTaskListPage />
          </ProtectedRoute>
        }
      />

      <Route
        path="/tasks/new"
        element={
          <ProtectedRoute>
            <TaskCreatePage />
          </ProtectedRoute>
        }
      />

      <Route
        path="/projects/:projectId/tasks/new"
        element={
          <ProtectedRoute>
            <TaskCreatePage />
          </ProtectedRoute>
        }
      />

      <Route
        path="/projects/:projectId/tasks/:taskId"
        element={
          <ProtectedRoute>
            <TaskDetailPage />
          </ProtectedRoute>
        }
      />

      <Route
        path="/projects/:projectId/tasks/:taskId/edit"
        element={
          <ProtectedRoute>
            <TaskEditPage />
          </ProtectedRoute>
        }
      />

      <Route
        path="*"
        element={
          <Navigate
            to="/"
            replace
          />
        }
      />
    </Routes>
  );
}

export default App;