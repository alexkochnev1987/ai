import React from "react";
import styles from "./Input.module.css";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  fullWidth?: boolean;
}

export const Input: React.FC<InputProps> = ({
  label,
  error,
  fullWidth = false,
  className = "",
  ...props
}) => {
  const containerClasses = [styles.container, fullWidth ? styles.fullWidth : ""]
    .filter(Boolean)
    .join(" ");

  const inputClasses = [
    styles.input,
    error ? styles.error : "",
    fullWidth ? styles.fullWidth : "",
    className,
  ]
    .filter(Boolean)
    .join(" ");

  return (
    <div className={containerClasses}>
      {label && <label className={styles.label}>{label}</label>}
      <input className={inputClasses} {...props} />
      {error && <p className={styles.errorMessage}>{error}</p>}
    </div>
  );
};
