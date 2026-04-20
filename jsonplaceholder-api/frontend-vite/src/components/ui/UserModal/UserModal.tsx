import React, { useEffect } from "react";
import { Button } from "../../common/Button";
import styles from "./UserModal.module.css";
import type { User } from "../../../types/user";

interface UserModalProps {
  user: User | null;
  isOpen: boolean;
  onClose: () => void;
}

export const UserModal: React.FC<UserModalProps> = ({
  user,
  isOpen,
  onClose,
}) => {
  useEffect(() => {
    const handleEscapeKey = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener("keydown", handleEscapeKey);
      document.body.style.overflow = "hidden";
    }

    return () => {
      document.removeEventListener("keydown", handleEscapeKey);
      document.body.style.overflow = "unset";
    };
  }, [isOpen, onClose]);

  if (!isOpen || !user) return null;

  const mapUrl = `https://maps.google.com/?q=${user.address.geo.lat},${user.address.geo.lng}`;

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div
        className={styles.modal}
        onClick={(e) => e.stopPropagation()}
        role="dialog"
        aria-labelledby="modal-title"
        aria-modal="true"
      >
        <div className={styles.modalHeader}>
          <h2 id="modal-title" className={styles.modalTitle}>
            {user.name}
          </h2>
          <Button
            variant="secondary"
            size="sm"
            className={styles.closeButton}
            onClick={onClose}
            aria-label="Close modal"
          >
            ×
          </Button>
        </div>

        <div className={styles.modalContent}>
          <div className={styles.section}>
            <h3>Contact Information</h3>
            <div className={styles.infoGrid}>
              <div className={styles.infoItem}>
                <span className={styles.label}>Username:</span>
                <span>{user.username}</span>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.label}>Email:</span>
                <a href={`mailto:${user.email}`}>{user.email}</a>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.label}>Phone:</span>
                <a href={`tel:${user.phone}`}>{user.phone}</a>
              </div>
              <div className={styles.infoItem}>
                <span className={styles.label}>Website:</span>
                <a
                  href={`https://${user.website}`}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {user.website}
                </a>
              </div>
            </div>
          </div>

          <div className={styles.section}>
            <h3>Address</h3>
            <div className={styles.addressInfo}>
              <p>
                {user.address.street}, {user.address.suite}
                <br />
                {user.address.city}, {user.address.zipcode}
              </p>
              <div className={styles.mapLink}>
                <a
                  href={mapUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className={styles.mapButton}
                >
                  📍 View on Map
                </a>
                <span className={styles.coordinates}>
                  Lat: {user.address.geo.lat}, Lng: {user.address.geo.lng}
                </span>
              </div>
            </div>
          </div>

          <div className={styles.section}>
            <h3>Company</h3>
            <div className={styles.companyInfo}>
              <h4>{user.company.name}</h4>
              <p className={styles.catchPhrase}>{user.company.catchPhrase}</p>
              <p className={styles.bs}>{user.company.bs}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default UserModal;
