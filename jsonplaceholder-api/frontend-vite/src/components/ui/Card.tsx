import React from "react";
import type { ReactNode } from "react";
import styles from "./Card.module.css";

interface CardProps {
  children: ReactNode;
  className?: string;
  padding?: "sm" | "md" | "lg";
}

export const Card: React.FC<CardProps> = ({
  children,
  className = "",
  padding = "md",
}) => {
  const classes = [styles.card, styles[padding], className]
    .filter(Boolean)
    .join(" ");

  return <div className={classes}>{children}</div>;
};
