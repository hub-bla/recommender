import pandas as pd

from sentence_transformers import SentenceTransformer

import psycopg2


def convert_titles(titles: pd.Series, model: SentenceTransformer) -> pd.DataFrame:
    """
        Create title embedding

        returns: DataFrame(title, embedding)
    """

    embeddings = model.encode(titles.to_numpy())
    return pd.concat([titles, pd.DataFrame(embeddings)], axis=1)


if __name__ ==  "__main__":
    model = SentenceTransformer('Mihaiii/gte-micro-v2')

    df = pd.read_csv('./movie_dataset.csv')

    title_embeddings = convert_titles(df['original_title'], model)

    connection = psycopg2.connect(database="movie_recommender", user="postgres", password="admin", host="localhost", port=5432)

    cursor = connection.cursor()

    cursor.execute("""
            CREATE EXTENSION vector;
            """)
    cursor.execute("""
                   CREATE TABLE title_embeddings (
                    id BIGSERIAL PRIMARY KEY,
                   title VARCHAR(100) NOT NULL,
                   embedding vector(384) NOT NULL
                   );
                """)
    
    for idx, row in title_embeddings.iterrows():
        my_doc = {"id": idx, 
                  "title": row.iloc[0], 
                  "embedding": row.iloc[1:].to_numpy().tolist()}
        query = """
        INSERT INTO title_embeddings(id, title, embedding) 
        VALUES (%(id)s, %(title)s, %(embedding)s)
        """
        cursor.execute(query, my_doc)
        
    connection.commit()

    cursor.execute('SELECT * FROM title_embeddings;')
    
    for r in cursor.fetchall(): # test if data was added
        print(r)
        break
    cursor.close()