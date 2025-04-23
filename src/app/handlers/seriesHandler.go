package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"seriesapp/src/app/models"

	"github.com/go-chi/chi/v5"

)

func GetSeriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var series []models.SeriesModel
		rows, err := db.Query("SELECT * FROM series")
		if err != nil {
			http.Error(w, "Error al consultar las series: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var s models.SeriesModel
			err := rows.Scan(
				&s.ID,
				&s.Title,
				&s.Status,
				&s.Episodes,
				&s.LastEpisode,
				&s.Ranking,
			)
			if err != nil {
				http.Error(w, "Error al leer los datos: "+err.Error(), http.StatusInternalServerError)
				return
			}
			series = append(series, s)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, "Error después de leer las series: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(series); err != nil {
			http.Error(w, "Error al generar la respuesta JSON", http.StatusInternalServerError)
		}
		
	}
}

func GetSeriesByIdHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")		
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID debe ser un número entero", http.StatusBadRequest)
			return
		}

		var serie models.SeriesModel
		err = db.QueryRow(
			"SELECT * FROM series WHERE id = ?",
			id,
		).Scan(&serie.ID, &serie.Title, &serie.Status, &serie.Episodes, &serie.LastEpisode, &serie.Ranking)

		switch {
		case err == sql.ErrNoRows:
			http.Error(w, "Serie no encontrada", http.StatusNotFound)
			return
		case err != nil:
			http.Error(w, "Error al consultar la serie: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(serie)
	}
}

func PostSeriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newSeries models.SeriesModel
		err := json.NewDecoder(r.Body).Decode(&newSeries)
		if err != nil {
			http.Error(w, "Error al decodificar el JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if newSeries.Title == "" {
			http.Error(w, "El título es requerido", http.StatusBadRequest)
			return
		}

		result, err := db.Exec(
			"INSERT INTO series (title, status, episodes, last_episode, ranking) VALUES (?, ?, ?, ?, ?)",
			newSeries.Title,
			newSeries.Status,
			newSeries.Episodes,
			newSeries.LastEpisode,
			newSeries.Ranking,
		)
		if err != nil {
			http.Error(w, "Error al insertar la serie: "+err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Error al obtener el ID: "+err.Error(), http.StatusInternalServerError)
			return
		}

		newSeries.ID = int(id)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newSeries)
	}
}

func PutSeriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		var updatedSeries models.SeriesModel
		err = json.NewDecoder(r.Body).Decode(&updatedSeries)
		if err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		if updatedSeries.Title == "" {
			http.Error(w, "El título es requerido", http.StatusBadRequest)
			return
		}

		result, err := db.Exec(
			"UPDATE series SET title = ?, status = ?, episodes = ? WHERE id = ?",
			updatedSeries.Title,
			updatedSeries.Status,
			updatedSeries.Episodes,
			id,
		)
		if err != nil {
			http.Error(w, "Error al actualizar la serie", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, "Error al verificar la actualización", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Serie no encontrada", http.StatusNotFound)
			return
		}
		updatedSeries.ID = id
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedSeries)
	}
}

func DeleteSeriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("DELETE FROM series WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Error al eliminar la serie", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, "Error al verificar la eliminación", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Serie no encontrada", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

//Handlers para los aspectos que valen punteo

func UpdateSeriesStatusHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "ID debe ser un número entero", http.StatusBadRequest)
            return
        }

        var updateData struct {
            Status string `json:"status"`
        }
        err = json.NewDecoder(r.Body).Decode(&updateData)
        if err != nil {
            http.Error(w, "Datos inválidos", http.StatusBadRequest)
            return
        }

        validStatus := map[string]bool{
            "Plan to Watch": true, 
			"Watching": true, 
            "Completed": true,
			"Dropped": true,
        }
        if !validStatus[updateData.Status] {
            http.Error(w, "Estado no válido", http.StatusBadRequest)
            return
        }

        var updatedSeries models.SeriesModel
        err = db.QueryRow(`
            UPDATE series SET status = ? 
            WHERE id = ? 
            RETURNING id, title, status, episodes, last_episode, ranking`,
            updateData.Status, id,
        ).Scan(
            &updatedSeries.ID, &updatedSeries.Title, &updatedSeries.Status,
            &updatedSeries.Episodes, &updatedSeries.LastEpisode, &updatedSeries.Ranking,
        )

        if err == sql.ErrNoRows {
            http.Error(w, "Serie no encontrada", http.StatusNotFound)
            return
        }
        if err != nil {
            http.Error(w, "Error al actualizar: "+err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(updatedSeries)
    }
}

func IncrementEpisodeHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "ID debe ser un número entero", http.StatusBadRequest)
            return
        }

        tx, err := db.Begin()
        if err != nil {
            http.Error(w, "Error al iniciar transacción", http.StatusInternalServerError)
            return
        }
        defer tx.Rollback()

        var totalEpisodes, lastEpisode int
        err = tx.QueryRow(`
            SELECT episodes, last_episode FROM series WHERE id = ?`,
            id,
        ).Scan(&totalEpisodes, &lastEpisode)
        
        if err == sql.ErrNoRows {
            http.Error(w, "Serie no encontrada", http.StatusNotFound)
            return
        }
        if err != nil {
            http.Error(w, "Error al consultar serie", http.StatusInternalServerError)
            return
        }

        if lastEpisode >= totalEpisodes {
            http.Error(w, "No hay más episodios disponibles", http.StatusBadRequest)
            return
        }

        var updatedSeries models.SeriesModel
        err = tx.QueryRow(`
            UPDATE series SET last_episode = last_episode + 1 
            WHERE id = ? 
            RETURNING id, title, status, episodes, last_episode, ranking`,
            id,
        ).Scan(
            &updatedSeries.ID, &updatedSeries.Title, &updatedSeries.Status,
            &updatedSeries.Episodes, &updatedSeries.LastEpisode, &updatedSeries.Ranking,
        )
        
        if err != nil {
            http.Error(w, "Error al actualizar", http.StatusInternalServerError)
            return
        }

        if err = tx.Commit(); err != nil {
            http.Error(w, "Error al confirmar cambios", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(updatedSeries)
    }
}

func UpvoteSeriesHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "ID debe ser un número entero", http.StatusBadRequest)
            return
        }

        var updatedSeries models.SeriesModel
        err = db.QueryRow(`
            UPDATE series SET ranking = ranking + 1 
            WHERE id = ? 
            RETURNING id, title, status, episodes, last_episode, ranking`,
            id,
        ).Scan(
            &updatedSeries.ID, &updatedSeries.Title, &updatedSeries.Status,
            &updatedSeries.Episodes, &updatedSeries.LastEpisode, &updatedSeries.Ranking,
        )

        if err == sql.ErrNoRows {
            http.Error(w, "Serie no encontrada", http.StatusNotFound)
            return
        }
        if err != nil {
            http.Error(w, "Error al actualizar ranking", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(updatedSeries)
    }
}

func DownvoteSeriesHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "ID debe ser un número entero", http.StatusBadRequest)
            return
        }

        var updatedSeries models.SeriesModel
        err = db.QueryRow(`
            UPDATE series SET ranking = CASE 
                WHEN ranking > 0 THEN ranking - 1 
                ELSE 0 
            END
            WHERE id = ? 
            RETURNING id, title, status, episodes, last_episode, ranking`,
            id,
        ).Scan(
            &updatedSeries.ID, &updatedSeries.Title, &updatedSeries.Status,
            &updatedSeries.Episodes, &updatedSeries.LastEpisode, &updatedSeries.Ranking,
        )

        if err == sql.ErrNoRows {
            http.Error(w, "Serie no encontrada", http.StatusNotFound)
            return
        }
        if err != nil {
            http.Error(w, "Error al actualizar ranking", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(updatedSeries)
    }
}

