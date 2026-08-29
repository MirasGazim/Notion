import { useState, FormEvent } from "react";
import { useNavigate, Link } from "react-router-dom";
import { signUp } from "../api/auth";
import { ApiError } from "../api/client";
import { useAuth } from "../auth/AuthContext";
import styles from "./SignUp.module.css";

interface FormErrors {
  email?: string;
  username?: string;
  password?: string;
}

export default function SignUp() {
  const navigate = useNavigate();
  const { login } = useAuth();

  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [errors, setErrors] = useState<FormErrors>({});
  const [serverError, setServerError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  function validate(): FormErrors {
    const next: FormErrors = {};
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      next.email = "Введите корректный email";
    }
    if (username.trim().length < 3) {
      next.username = "Минимум 3 символа";
    }
    if (password.length < 6) {
      next.password = "Минимум 6 символов";
    }
    return next;
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setServerError(null);

    const validationErrors = validate();
    setErrors(validationErrors);
    if (Object.keys(validationErrors).length > 0) return;

    setLoading(true);
    try {
      const { token } = await signUp({ email, username, password });
      login(token);
      navigate("/");
    } catch (err) {
      if (err instanceof ApiError && err.status === 409) {
        setServerError("Пользователь с таким email или username уже существует");
      } else {
        setServerError("Что-то пошло не так, попробуйте позже");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className={styles.wrapper}>
      <form className={styles.card} onSubmit={handleSubmit} noValidate>
        <h1 className={styles.title}>Создать аккаунт</h1>

        <label className={styles.field}>
          <span>Email</span>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className={errors.email ? styles.inputError : undefined}
          />
          {errors.email && <span className={styles.error}>{errors.email}</span>}
        </label>

        <label className={styles.field}>
          <span>Username</span>
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className={errors.username ? styles.inputError : undefined}
          />
          {errors.username && <span className={styles.error}>{errors.username}</span>}
        </label>

        <label className={styles.field}>
          <span>Пароль</span>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className={errors.password ? styles.inputError : undefined}
          />
          {errors.password && <span className={styles.error}>{errors.password}</span>}
        </label>

        {serverError && <div className={styles.serverError}>{serverError}</div>}

        <button type="submit" className={styles.submit} disabled={loading}>
          {loading ? "Создаём..." : "Зарегистрироваться"}
        </button>

        <p className={styles.switch}>
          Уже есть аккаунт? <Link to="/sign-in">Войти</Link>
        </p>
      </form>
    </div>
  );
}