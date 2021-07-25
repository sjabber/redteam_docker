package com.hanium.mer;

import io.jsonwebtoken.*;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import javax.servlet.http.Cookie;
import java.io.UnsupportedEncodingException;

@Service
public class TokenUtils {

    @Value("${Key}")
    private String KEY;
    private static final int SUCCESS = 0;
    private static final int EXPIRE = 1;
    private static final int TAMPER = 2;
    private static final int NULL = 3;

    public int isValidToken(String token) throws UnsupportedEncodingException, ExpiredJwtException, JwtException, NullPointerException  {
        Claims claims = Jwts.parser().setSigningKey(KEY.getBytes("UTF-8")).parseClaimsJws(token).getBody();

        return SUCCESS;
    }

    public Claims getClaimsFormToken(Cookie[] cookies) throws UnsupportedEncodingException {

        Claims claims = null;

        for(Cookie c : cookies) {
            if (c.getName().equals("access-token")) {
                claims = Jwts.parser().setSigningKey(KEY.getBytes("UTF-8")).parseClaimsJws((c.getValue())).getBody();
            }
        }
        return claims;
    }


    public <T> String create() throws UnsupportedEncodingException{
        String jwt = Jwts.builder()
                .setHeaderParam("alg","HS256")
                .setHeaderParam("typ", "JWT")
                .claim("authorized", true)
                .claim("exp", 1610889330)
                .claim("user_email","gc5215@nate.com")
                .claim("user_name", "Kwangchae Kim")
                .claim("user_no", 4)
                .signWith(SignatureAlgorithm.HS256, KEY.getBytes("UTF-8"))  //같은 값나옴
                .compact();
        return jwt;
    }
}

